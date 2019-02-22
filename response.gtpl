<html>
  <head>
         <title>Magnificent image classifier</title>
         <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
         <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css" integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous">
  </head>
  <body>
  <div  id="app">
    <nav class="navbar navbar-expand-lg navbar-light" style="background-color: #e3f2fd;">
            <div class="container-fluid">
                <div class="navbar-header">
                    <a class="navbar-brand">Magnificent Image Classifier</a>
                </div>
                
            </div>
    </nav>

    <div class="container" style="margin-top: 10px;">
    <div class="card" >
          <div class="card-header">
              Inference Result
          </div>
          <div class="card-body">
              Category: {{.Category}}<BR>
              Probability: {{.Probability}}<br></br>
              <img src="data:image/jpg;base64,{{.Picture}}"><br><br>
              <a class="btn btn-primary btn-sm" href="/">Go Back</a>
            </div>
          
          </div>
    </div>
    </div>
  </body>
  </html>