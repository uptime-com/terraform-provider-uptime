# Import using the credential ID
# The secret itself is write-only - the API never returns it - so the first plan after
# importing proposes setting the secret from your configuration.
terraform import uptime_credential.example 123
