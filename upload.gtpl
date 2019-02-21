<html>
  <head>
         <title>Magnificent image classifier</title>
  </head>
  <body>
  <h3>Upload an image and enjoy the magic...</h3>
<form enctype="multipart/form-data" action="http://127.0.0.1:9090/" method="post">
          {{/* 1. File input */}}
          <input type="file" name="uploadfile" />
 
          {{/* 2. Submit button */}}
          <input type="submit" value="upload file" />
      </form>
 
  </body>
  </html>