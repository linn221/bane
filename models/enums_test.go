package models

import (
	"bytes"
	"testing"
	"time"
)

func TestMyDate_UnmarshalGQL(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name      string
		input     string
		wantYear  int
		wantMonth time.Month
		wantDay   int
		wantErr   bool
	}{
		{
			name:      "jan 8 without year",
			input:     "jan 8",
			wantYear:  currentYear,
			wantMonth: time.January,
			wantDay:   8,
			wantErr:   false,
		},
		{
			name:      "jan 8 2024 with year",
			input:     "jan 8 2024",
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   8,
			wantErr:   false,
		},
		{
			name:      "feb 15 2023",
			input:     "feb 15 2023",
			wantYear:  2023,
			wantMonth: time.February,
			wantDay:   15,
			wantErr:   false,
		},
		{
			name:      "dec 25 without year",
			input:     "dec 25",
			wantYear:  currentYear,
			wantMonth: time.December,
			wantDay:   25,
			wantErr:   false,
		},
		{
			name:      "sept 1 2025",
			input:     "sept 1 2025",
			wantYear:  2025,
			wantMonth: time.September,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:      "sep 1 2025",
			input:     "sep 1 2025",
			wantYear:  2025,
			wantMonth: time.September,
			wantDay:   1,
			wantErr:   false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid month",
			input:   "xyz 8",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d MyDate
			err := d.UnmarshalGQL(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalGQL() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if d.Time.Year() != tt.wantYear {
				t.Errorf("Year() = %v, want %v", d.Time.Year(), tt.wantYear)
			}
			if d.Time.Month() != tt.wantMonth {
				t.Errorf("Month() = %v, want %v", d.Time.Month(), tt.wantMonth)
			}
			if d.Time.Day() != tt.wantDay {
				t.Errorf("Day() = %v, want %v", d.Time.Day(), tt.wantDay)
			}
		})
	}
}

func TestMyDate_MarshalGQL(t *testing.T) {
	tests := []struct {
		name string
		date MyDate
		want string
	}{
		{
			name: "January 8 2025",
			date: MyDate{Time: time.Date(2025, time.January, 8, 0, 0, 0, 0, time.UTC)},
			want: `"January 8 2025"`,
		},
		{
			name: "December 25 2024",
			date: MyDate{Time: time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)},
			want: `"December 25 2024"`,
		},
		{
			name: "February 15 2023",
			date: MyDate{Time: time.Date(2023, time.February, 15, 0, 0, 0, 0, time.UTC)},
			want: `"February 15 2023"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.date.MarshalGQL(&buf)
			got := buf.String()

			if got != tt.want {
				t.Errorf("MarshalGQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyTime_UnmarshalGQL(t *testing.T) {
	currentYear := time.Now().Year()

	tests := []struct {
		name       string
		input      string
		wantYear   int
		wantMonth  time.Month
		wantDay    int
		wantHour   int
		wantMinute int
		wantErr    bool
	}{
		{
			name:       "jan 8 11:01 PM without year",
			input:      "jan 8 11:01 PM",
			wantYear:   currentYear,
			wantMonth:  time.January,
			wantDay:    8,
			wantHour:   23,
			wantMinute: 1,
			wantErr:    false,
		},
		{
			name:       "jan 8 2024 11:01 PM with year",
			input:      "jan 8 2024 11:01 PM",
			wantYear:   2024,
			wantMonth:  time.January,
			wantDay:    8,
			wantHour:   23,
			wantMinute: 1,
			wantErr:    false,
		},
		{
			name:       "feb 15 9:30 AM",
			input:      "feb 15 9:30 AM",
			wantYear:   currentYear,
			wantMonth:  time.February,
			wantDay:    15,
			wantHour:   9,
			wantMinute: 30,
			wantErr:    false,
		},
		{
			name:       "dec 25 12:00 PM",
			input:      "dec 25 12:00 PM",
			wantYear:   currentYear,
			wantMonth:  time.December,
			wantDay:    25,
			wantHour:   12,
			wantMinute: 0,
			wantErr:    false,
		},
		{
			name:       "jan 1 12:00 AM (midnight)",
			input:      "jan 1 12:00 AM",
			wantYear:   currentYear,
			wantMonth:  time.January,
			wantDay:    1,
			wantHour:   0,
			wantMinute: 0,
			wantErr:    false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "missing AM/PM",
			input:   "jan 8 11:01",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mt MyTime
			err := mt.UnmarshalGQL(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UnmarshalGQL() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if mt.Time.Year() != tt.wantYear {
				t.Errorf("Year() = %v, want %v", mt.Time.Year(), tt.wantYear)
			}
			if mt.Time.Month() != tt.wantMonth {
				t.Errorf("Month() = %v, want %v", mt.Time.Month(), tt.wantMonth)
			}
			if mt.Time.Day() != tt.wantDay {
				t.Errorf("Day() = %v, want %v", mt.Time.Day(), tt.wantDay)
			}
			if mt.Time.Hour() != tt.wantHour {
				t.Errorf("Hour() = %v, want %v", mt.Time.Hour(), tt.wantHour)
			}
			if mt.Time.Minute() != tt.wantMinute {
				t.Errorf("Minute() = %v, want %v", mt.Time.Minute(), tt.wantMinute)
			}
		})
	}
}

func TestMyTime_MarshalGQL(t *testing.T) {
	tests := []struct {
		name string
		time MyTime
		want string
	}{
		{
			name: "January 8 11:01 PM",
			time: MyTime{Time: time.Date(2025, time.January, 8, 23, 1, 0, 0, time.UTC)},
			want: `"January 8 11:01 PM"`,
		},
		{
			name: "December 25 9:30 AM",
			time: MyTime{Time: time.Date(2024, time.December, 25, 9, 30, 0, 0, time.UTC)},
			want: `"December 25 9:30 AM"`,
		},
		{
			name: "February 15 12:00 PM",
			time: MyTime{Time: time.Date(2023, time.February, 15, 12, 0, 0, 0, time.UTC)},
			want: `"February 15 12:00 PM"`,
		},
		{
			name: "January 1 12:00 AM",
			time: MyTime{Time: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want: `"January 1 12:00 AM"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.time.MarshalGQL(&buf)
			got := buf.String()

			if got != tt.want {
				t.Errorf("MarshalGQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
