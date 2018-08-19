//XmlHttpRequest对象      
function createXmlHttpRequest(){      
    if(window.ActiveXObject){ //如果是IE浏览器      
        return new ActiveXObject("Microsoft.XMLHTTP");      
    }else if(window.XMLHttpRequest){ //非IE浏览器      
        return new XMLHttpRequest();      
    }      
} 


function build(){      
      	  
	//var url = "LoginServlet?username="+userName+"&password="+passWord+"";         
	var url = 'http://localhost:8081';          
	      
	xmlHttpRequest = createXmlHttpRequest();      
	     
	xmlHttpRequest.onreadystatechange = statechanged;      
		       
	xmlHttpRequest.open("POST",url,true);
	// xmlHttpRequest.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;");      
    
	xmlHttpRequest.send('{"BPM":70}');
	
}



// 回调函数      
function statechanged(){
	var req = xmlHttpRequest;
	if(req.readyState == 4 ){
	  if(req.status == 201){
			var json_str = xmlHttpRequest.responseText; // json形式的字符串
			questions = eval('(' + json_str + ')'); // 转化为json格式
			console.log(json_str)
			// var user = JSON.parse(json_str); // 转化为json格式的另一种方式，较安全
	  }
	  else if(req.status == 404)
	  {
			alert("request url is not found");
	  }
	  else if(req.status == 401 || req.status == 403)
	  {
			alert("request url is forbidden or not authorized to visit.");
	  }
	  else
	  {
			alert("unexpected error!Status Code :"+req.status);
	  }
	}                    
}

var btn = document.getElementById("a1");
console.log(btn)
btn.addEventListener('click', build)