# Import using composite ID: service_id:variable_id
# service_id is required because the API does not return it; importing by the
# variable ID alone would leave it unset and the next plan would propose a replacement.
# The provider verifies that service_id names the check that owns the variable and
# refuses the import otherwise.
terraform import uptime_service_variable.token 5891524:41218
