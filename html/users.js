function userFunction(org, group) {
  // Create a request to use to call the organization/org/users endpoint
  var request = new XMLHttpRequest()
  request.open('GET', 'https://y0319t11923.nordstrom.net:8111/organizations/'+org+'/groups/'+group, true)
  request.onload = function () {
    var data = JSON.parse(this.response)
    var users = JSON.parse(data)
    users.Array.sort();
  
    var html = "<table border='1|1'>";
    html+="<tr>"+group+"<tr>"
    for (var i = 0; i < users.Array.length; i++) {
      html+="<tr>";
      html+="<td class='userList' id=user"+org+group+users.Array[i]+">"+users.Array[i]+"</td>";
      html+="</tr>";
    }
    html+="</table>";
    console.log("GROUP"+group+org)
    document.getElementById(id="group"+org+group).innerHTML = html;
  }
  request.send()
}
