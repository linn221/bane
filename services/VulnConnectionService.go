package services

import (
	"fmt"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type vulnConnectionService struct {
	db *gorm.DB
}

// GetVulnsByReference retrieves all Vulns connected to a specific reference
func (s *vulnConnectionService) GetVulnsByReference(referenceId int, referenceType models.VulnReferenceType) ([]*models.Vuln, error) {
	var vulns []*models.Vuln

	err := s.db.
		Table("vulns").
		Joins("INNER JOIN vuln_connections ON vulns.id = vuln_connections.vuln_id").
		Where("vuln_connections.reference_id = ? AND vuln_connections.reference_type = ?", referenceId, referenceType).
		Find(&vulns).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get vulns for reference %d of type %s: %w", referenceId, referenceType, err)
	}

	return vulns, nil
}

// ConnectReferencesToVuln connects multiple vulns to a reference by their aliases
func (s *vulnConnectionService) ConnectReferencesToVuln(referenceType models.VulnReferenceType, referenceId int, vulnAliases []string) error {
	if len(vulnAliases) == 0 {
		return nil // Nothing to connect
	}

	// First, get the vuln IDs for the given aliases
	var vulnIds []int
	err := s.db.Model(&models.Vuln{}).
		Select("id").
		Where("alias IN ?", vulnAliases).
		Pluck("id", &vulnIds).Error

	if err != nil {
		return fmt.Errorf("failed to get vuln IDs for aliases: %w", err)
	}

	if len(vulnIds) == 0 {
		return fmt.Errorf("no vulns found for the provided aliases")
	}

	// Check if any aliases were not found
	if len(vulnIds) != len(vulnAliases) {
		return fmt.Errorf("some vuln aliases were not found")
	}

	// Create the connections
	var connections []models.VulnConnection
	for _, vulnId := range vulnIds {
		connections = append(connections, models.VulnConnection{
			VulnId:        vulnId,
			ReferenceId:   referenceId,
			ReferenceType: referenceType,
		})
	}

	// Use transaction to ensure atomicity
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First, remove existing connections for this reference
		err := tx.Where("reference_id = ? AND reference_type = ?", referenceId, referenceType).
			Delete(&models.VulnConnection{}).Error
		if err != nil {
			return fmt.Errorf("failed to remove existing connections: %w", err)
		}

		// Then create new connections
		if len(connections) > 0 {
			err = tx.Create(&connections).Error
			if err != nil {
				return fmt.Errorf("failed to create new connections: %w", err)
			}
		}

		return nil
	})
}

// DisconnectAllReferences removes all connections for a specific reference
func (s *vulnConnectionService) DisconnectAllReferences(referenceType models.VulnReferenceType, referenceId int) error {
	err := s.db.Where("reference_id = ? AND reference_type = ?", referenceId, referenceType).
		Delete(&models.VulnConnection{}).Error

	if err != nil {
		return fmt.Errorf("failed to disconnect references for %d of type %s: %w", referenceId, referenceType, err)
	}

	return nil
}

// GetReferenceCount returns the number of references connected to a specific vuln
func (s *vulnConnectionService) GetReferenceCount(vulnId int) (int64, error) {
	var count int64
	err := s.db.Model(&models.VulnConnection{}).
		Where("vuln_id = ?", vulnId).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to get reference count for vuln %d: %w", vulnId, err)
	}

	return count, nil
}

// GetVulnReferences returns all references connected to a specific vuln
func (s *vulnConnectionService) GetVulnReferences(vulnId int) ([]*models.VulnConnection, error) {
	var connections []*models.VulnConnection

	err := s.db.Where("vuln_id = ?", vulnId).Find(&connections).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get references for vuln %d: %w", vulnId, err)
	}

	return connections, nil
}
