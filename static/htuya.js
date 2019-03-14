function swSet(sw,val){
   // Feead back provided by data update 
   $.get("set",{'switch':sw,set:val});
}

$(document).ready(function(){
   dataUpdate(false)
});

function dataUpdate(w){
   // update switch status
   // the server blocks for 15 seconds unless data is updated
   $.post('get',{w:w},gotData,'json');
}
function gotData(data){
  var d = new Date();
  $.each( data, function( key, val ) {
    console.log(key+" " + val.Value + " " +val.Status);
    $("#sw-"+key).prop("checked",val.Value).change();
  });
  $("#wupdate").text("Update:"+d.toLocaleTimeString());
  dataUpdate(true);// loop
};
