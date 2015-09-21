//exec-util.js

var exec = require('child_process').exec;
module.exports={
	execCmd: function(cmd, callback){
		exec(cmd, function(error, stdout, stderr){
			callback(error, stdout, stderr);
		});
	},
	execCmd2: function(cmd, options, callback){
		exec(cmd, options, function(error, stdout, stderr){
			callback(error, stdout, stderr);
		});
	}
};