# Why POSTing 2 files instead of the archive as file and a json in the request body ?

Because HTTP protocol cannot handle a `multipart/form-data` request with a `body`.