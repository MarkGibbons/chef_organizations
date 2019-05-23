function organizationFunction() {
  // Create a request to use to call the organizations api
  var request = new XMLHttpRequest()
  // request.open('GET', 'https://m00973810.nordstrom.net:8080/organizations', true)
  // https://www.taniarascia.com/how-to-connect-to-an-api-with-javascript/
  request.open('GET', 'https://y0319t11923.nordstrom.net:8111/organizations', true)
  request.onload = function () {
    // Access JSON here
    console.log(this.response)
    console.log(typeof this.response)
    var data = JSON.parse(this.response)
    var orgs = JSON.parse(data)
    orgs.Array.sort();
    console.log(orgs.Arraylength)
    console.log(orgs.Array[1])
  
    var html = "<table border='1|1'>";
    for (var i = 0; i < orgs.Array.length; i++) {
      html+="<tr>";
      html+="<td class='orgList' id=organization"+orgs.Array[i]+" onclick=groupFunction('"+orgs.Array[i]+"')><u>"+orgs.Array[i]+"</u></td>";
      html+="</tr>";
    }
    html+="</table>";
    console.log(html)
    document.getElementById('organizationList').innerHTML = html;
  }
  request.send()
}
