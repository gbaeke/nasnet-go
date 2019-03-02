# storage account and share

Create the storage account

az storage account create --resource-group scratch-rg --location westeurope --sku Standard_LRS --name nasnetstorage

Create the share 

az storage share create --name certs --account-name nasnetstorage

Obtain the storage key

STORAGE_KEY=$(az storage account keys list --resource-group scratch-rg --account-name nasnetstorage --query "[0].value" --output tsv)
echo $STORAGE_KEY

The storage key needs to be entered in aci.yaml (search ACCOUNTKEY)