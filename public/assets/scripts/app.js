$(document).ready(function() {
    
    activateMethodLinks();
    
    activateFilterFields();
    
    
    activateShowlinks();
});


function activateMethodLinks() {
        
   
    // links with method post should post to the url at href instead of sending a get request
    $("a[method='post'], a[method='delete']").click(function(e){
        
        
        // Check on deletes whether the user really wants to do this
        if ($(this).attr('method') == 'delete') {
            if (!confirm('Are you sure you want to delete this item, this action cannot be undone?')) { 
                e.preventDefault();
                return false
            }
        } 
        

        
        
       var redirect = $(this).attr('data-redirect')
       var loc =  window.location
       
        // Perform a post to the href url
		$.ajax({
            type: "POST",
         url: $(this).attr('href').toString(),
            success: function(){
                if (redirect !== undefined && redirect.length > 0) {
                    window.location = redirect
                } else {
                    // Else just reload the page (lazy) -
                    //  we should be inserting fragments here
                    window.location.reload()
                }
            }
        })
        
        
       e.preventDefault();
       return false
    });
    
}

function activateFilterFields() {
	if ($('.filter-form').length == 0)  {
	    return
	}
   
    $('.filter-form .field select').change(function(e){
        this.form.submit();
    });
   

}


function activateShowlinks() {
    $(".show-link").click(function(e){
        $(this).parent().find(".show-summary").toggle()
        e.preventDefault()
    })
}