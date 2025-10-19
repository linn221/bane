add curl field in Endpoint type (forced resolver), taking KVPairGroup of variables as input for that field.
you may use a well maintained popular library for generating the curl command or create your own if it is relatively easy to do so.
when variables input is not empty, you will execute the related VarString with the variable. input will come in string "userId:2 category_id:33", and respective VarString type will inject the variables to its internal map, after that, generate the curl commands with variables plugged. Makes sure the command is valid, well escaped, etc.
