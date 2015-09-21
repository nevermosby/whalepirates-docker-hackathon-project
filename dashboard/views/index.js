/*
	page event for index.jade
*/
console.log($.fn.jquery);
var rootKey='swift-config';

$('a#backup').click(function(event){
	console.log('===start to backup===');
	alert('Your configuration has been successfully backuped');
/*	
	$.ajax({
		type : "POST",
		url: "/backup",
		dataType: 'json',
		data: ''
	})
	.done(function(result){
		if(result){
			alert('Your configuration has been successfully backuped');
		}else{
			alert('Server error');
		}

	})
	.fail(function(){
		alert('error')
	})
	.always(function(){
		console.log('backup end')
	})
*/
	event.preventDefault();
})

$('a#restore').click(function(event){
	console.log('====start to restore==');
	var selectedVersion = $('#restore-version option:selected').val();
	alert('Your configuration has been restored to ' + selectedVersion);
	event.preventDefault();
})

$('a#save').click(function(){
	console.log('===start to save===');

	var postData={};
	// loop the input list to form the target 
	$('.container').find('input').each(function(idx,ele){
		var $current = $(ele)
		var key='', val='';
		//console.log($current);
		//console.log($current.attr('id'));
		key=$current.attr('id');
		// console.log($current.text());
		if ($current.attr('type') === 'checkbox'){
			//console.log($current.is(':checked'))
			val=$current.is(':checked');
		} else {
			//console.log($current.val())
			val=$current.val();
		}
		postData[key]=val;
	});
	console.log(JSON.stringify(postData));
	
	// post the data to backend
	$.ajax({
		type: "POST",
		url: "/save",
		dataType: 'json',
		data: JSON.stringify(postData)
	})
	.done(function(result){
		console.log("success return: " + result);
                alert('Your configuraton has been saved.')
	})
	.fail(function(){
		alert('error');
	})
	.always(function(){
		console.log('always');
	});
});

