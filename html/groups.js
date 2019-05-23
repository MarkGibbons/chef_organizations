function groupFunction(org) {
  // Create a request to use to call the organization/org/groups endpoint
  var request = new XMLHttpRequest()
  request.open('GET', 'https://y0319t11923.nordstrom.net:8111/organizations/'+org+'/groups', true)
  request.onload = function () {
    console.log(this.response)
    console.log(typeof this.response)
    var data = JSON.parse(this.response)
    var groups = JSON.parse(data)
    groups.Array.sort();
    console.log(groups.Arraylength)
    console.log(groups.Array[1])
  
    var html = "<table border='1|1'>";
    html+="<tr>"+org+"<tr>"
    for (var i = 0; i < groups.Array.length; i++) {
      html+="<tr>";
      html+="<td class='groupList' id=group"+org+groups.Array[i]+" onclick=userFunction('"+org+"','"+groups.Array[i]+"')><u>"+groups.Array[i]+"</u></td>";
      // html+="<td id=group"+groups.Array[i]+">"+groups.Array[i]+"</td>";
      html+="</tr>";
    }
    html+="</table>";
    console.log(html)
    document.getElementById(id="organization"+org).innerHTML = html;
  }
  request.send()
}
