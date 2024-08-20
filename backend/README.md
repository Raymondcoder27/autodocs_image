# How to Use

**NOTE:** The exposed port to localhost is `1001`.

## Endpoints

The program has 5 endpoints:

1. **Upload Template**
   - **Endpoint:** `POST /upload-template`
   - **Handler:** `controllers.UploadTemplate`
   - **Description:** Takes an HTML template file and returns a reference number in the format "D960888-0002" for later use in PDF generation.

2. **Generate Document**
   - **Endpoint:** `POST /generate/:id`
   - **Handler:** `controllers.CreateDocument`
   - **Description:** Generates a PDF document using the provided reference number. Returns the PDF as a base64-encoded string.

3. **Get All Documents**
   - **Endpoint:** `GET /documents`
   - **Handler:** `controllers.GetDocuments`
   - **Description:** Retrieves all document details, including the original file name, MinIO template file name, and reference number.

4. **Get Document by ID**
   - **Endpoint:** `GET /documents/:id`
   - **Handler:** `controllers.GetDocumentById`
   - **Description:** Retrieves a specific PDF by its ID. Returns the PDF as a base64-encoded string, along with other details.

5. **Delete Document**
   - **Endpoint:** `DELETE /documents/:id`
   - **Handler:** `controllers.DeleteDocument`
   - **Description:** Deletes the PDF file from MinIO and removes its details from the database.

## Example Usage

```bash
# Uploading a Template File
curl -X POST http://localhost:1001/upload-template \
     -F "template=@/path/to/your/template.html"

# Generating a Document
curl -X POST http://localhost:1001/generate/document_id \
     -F "ref_no=DXXXXXX-XXXX" \
     -F "data={\"key\":\"value\"}"

# Retrieving All Documents
curl -X GET http://localhost:1001/documents

# Retrieving a Specific Document by ID
curl -X GET http://localhost:1001/documents/your_document_id

# Deleting a Document
curl -X DELETE http://localhost:1001/documents/your_document_id
