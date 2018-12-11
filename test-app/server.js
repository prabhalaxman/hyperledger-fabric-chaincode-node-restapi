var express = require('Express');
var app = express();
var api = require('./api/blockChainApi');

app.use('/api', api);

app.get('/', function(req, res) {
    res.send('Hello! This is the sample REST API to access Hyperledger Blockchain smartcontract.The API is at http://localhost:' + port + '/api');
});

 var port = process.env.PORT || 8080; // used to create, sign, and verify tokens
 app.listen(port);