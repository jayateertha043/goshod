function shodanSearch(){

    searchbar = document.getElementById("search").value;
    search = document.getElementById("search").value;
    console.log("Executing shodanSearch ...")
    console.log(searchbar)
 

    var xhttp = new XMLHttpRequest();
    var url = '/shodansearch';

    xhttp.open("POST", url, true);
    xhttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    var response_body="";
    xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
           // Typical action to be performed when the document is ready:
           response_body = xhttp.responseText
           console.log(response_body);
        }
    };
    xhttp.send("search="+search);

    
}