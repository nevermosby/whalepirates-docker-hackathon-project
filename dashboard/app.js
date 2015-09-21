var express = require('express');
var path = require('path');

// body parsing setup
var bodyParser = require('body-parser');

var fs = require('fs');
var request = require('request');
var Etcd = require('node-etcd');
etcd = new Etcd("ETCD_SERVER_ADDRESS", "ETCD_SERVER_PORT");
var execUtil = require('./exec-util');

// Constants
var PORT = 4000;

// App
var app = express();

app.use(bodyParser.json()); // for parsing application/json
app.use(bodyParser.urlencoded({ extended : true })); // for parsing application/x-www-form-urlencoded

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');
app.use(express.static(path.join(__dirname, 'public')));

app.get('/', function (req, res) {
	//res.send('Hello world\n');
	res.render('index', {
		title : "Etcd Configuration Dashboard"
	})
});

/*
app.post('/backup', function(req, res){

	var cmd ='./backup.sh';
	execUtil.execCmd(cmd, function(error, stdout, stderr){
		if(error){                                     
			console.error('error running: ' + cmd);
			console.error(stderr);
			res.json('');                          
		}else{                                         
        		result = stdout.split("\n");           
        		res.json(result);                      
		}                                              
	});

});
*/

app.post('/save', function (req, res) {
	configBody = req.body;
	configBodyStr = '';
	if (configBody) {
		configBodyStr = JSON.stringify(configBody);
	}
	//console.log("input raw string: " + configBodyStr);
	// handle bug introduced by posted data with additional colon
	tmpKeys = Object.keys(configBody)
		console.log("key list:" + tmpKeys);
	configBody = JSON.parse(tmpKeys);
	for (var p in configBody) {
		console.log('key:' + p);
		console.log('val:' + configBody[p]);
	}
	// save configuration data into etcd server
	etcd.raw('PUT', 'v2/keys/hack/storage/swift', configBody, {}, function (error, result) {
		if (error) {
			console.log('Put error with etcd server: ' + error);
			res.json('error');
		} else {
			console.log(result);
			res.json('OK');
			}
	});
});

app.listen(PORT);
console.log('Running on http://localhost:' + PORT);
