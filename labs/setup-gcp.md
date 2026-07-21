edit_file "setup-gcp.md" {
  instructions = "I will create a new setup file for GCP configuration"
  code_edit = """# Setup Student VM's

# Open the lab repo in VS Code

### Step 1: Open the folder

1. Open Visual Studio Code
2. In Visual Studio Code, click **File** -> **Open Folder** and browse to `C:\Users\TEKstudent\Downloads\repos\tf-dev`
3. After opening the folder, click the third icon in the left toolbar for source control. Next to **changes**, click the three dots and choose **pull**.

# Configure GCP credentials

### **Step 1: Set up GCP Service Account**

1. In a browser, log into the [Google Cloud Console](https://console.cloud.google.com/) using the credentials in the spreadsheet below.
    * [Cloud credentials](https://docs.google.com/spreadsheets/d/1gTV6btPeIyyXylRkDn2_LNbWkf9BGU6wsi5eIb-ynLY/edit?gid=2103659978#gid=2103659978)
2. Search for "IAM & Admin" in the search bar.
3. Click **Service Accounts**
4. Click **CREATE SERVICE ACCOUNT**
5. Name it "terraform-admin"
6. Click **CREATE AND CONTINUE**
7. Assign the "Editor" role
8. Click **DONE**
9. Click on the service account you just created
10. Click **KEYS** tab
11. Click **ADD KEY** -> **Create new key**
12. Choose **JSON** format
13. **IMPORTANT:** The key file will automatically download. Keep this file secure.

---

### **Step 2: Configure GCloud CLI**

1. In the Visual Studio Code Terminal run:

   ```
   gcloud auth activate-service-account --key-file=PATH_TO_YOUR_SERVICE_ACCOUNT_KEY.json
   ```

2. Set the default project:
   ```
   gcloud config set project PROJECT_ID
   ```
3. Set the default region:
   ```
   gcloud config set compute/region us-west1
   ```

---

### **Step 3: Test GCloud CLI**

1. Confirm that `gcloud` is properly authenticated:

   ```
   gcloud auth list
   ```

2. You should see your service account listed as active.

## Congratulations

Congratulations, you've successfully configured your machine.