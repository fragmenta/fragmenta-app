$(document).ready(function() {
    activateEditableToolbar()
})

// TODO - remove dependency on jquery?
// TODO - show formatting selected on toolbar when selection changes in our contenteditable
// TODO - Intercept paste on editable and remove complex html before it is pasted (from word etc)
// TODO - perhaps intercept return key to be sure we get a para, and to be sure we insert newline within code sections, not new code tags + br or similar
      
      
function activateEditableToolbar() {
    var toolbars = $('.content-editable-toolbar')
	if (toolbars.length == 0) {
        return;
    }
    
    toolbars.each(function(){
        var t = $(this)
        t.editable = $('#' + t.attr('data-editable') + '-editable' )
        t.textarea = $('#' + t.attr('data-editable') + '-textarea' )
        
        // Listen to a form submit, and call updateContent to make sure 
        // our textarea in the form is up to date with the latest content
       
       $(t.textarea[0].form).submit(function(e) {
           updateContent(t.editable,t.textarea,false)
       })


        // Intercept paste on editable and remove complex html before it is pasted here?
       t.editable.on('input',function(e) {
          cleanHTMLElements(t.editable)
       })
        
        
        // Listen to button clicks
        t.buttons = t.find('a')
       
        t.buttons.each(function(){
            // Link this button to a command
            var b = $(this)
            b.click(function(e){
               var cmd = this.id
               var insert = ""
                
                switch (cmd){
                case "showCode":
                    updateContent(t.editable,t.textarea,true)
                break;
                case "createLink": 
                    insert = prompt("Supply the web URL to link to");
                    if (!startsWith(insert,'http')) {
                        insert = "http://" + insert
                    }
                break;
                case "formatblock": 
                    insert = $(this).attr('data-format');
                break;
                default:
                break;
                }
                
               
               if (cmd.length > 0) {
                    document.execCommand(cmd,false,insert)
               }
                
            
               // Find and remove evil html created by browsers 
               
               
             
               var sel = $(getSelectionParentElement())
               
            
               if (sel.length > 0) {
                  // Clean align stuff
                  cleanAlign(cmd,sel)
                  cleanHTMLElements(sel)
                  sel.removeAttr('style') 
               } 
               
    
                e.preventDefault()
            })
            
        })
        
        
        
    })
    
    
   
    

}


// If textarea visible update the content editable with new html
// else update the textarea with new content editable html
function updateContent(editable,textarea,toggle) {
    var html = ""
    
    if (textarea.is(':visible')) {
        html = textarea.val()
        editable[0].innerHTML = html
        if (toggle){
            editable.show()
            textarea.hide()
        }
    } else {
        html = editable[0].innerHTML
        // Cleanup the html by removing plain spans
        html = cleanHTML(html)
        textarea.val(html)
        if (toggle){
            editable.hide()
            textarea.show()
        }
    }


console.log("Updating content"+html)
}

// Purge html string of plain spans and nbsp
function cleanHTML(html) {
    html = html.replace(/<\/?span>/gi,'')// Remove all empty span tags
    html = html.replace(/<\/?font [^>]*>/gi,'')// Remove ALL font tags
    html = html.replace(/&nbsp;</gi,' <')// this is sometimes required, be careful
    html = html.replace(/<p><\/p>/gi,'\n')// pretty format but remove empty paras
    html = html.replace(/<br><\/li>/gi,'<\/li>')
    
    // Remove comments and other MS cruft
    html = html.replace(/<!--[\w\d\[\]\s<\/>:.!="*{};-]*-->/gi,'')
    html = html.replace(/ class\=\"MsoNormal\"/gi,'')
    html = html.replace(/<p><o:p> <\/o:p><\/p>/gi,'')
    
    
    
    // Pretty printing elements which follow on from one another
    html = html.replace(/><(li|ul|ol|p|h\d|\/ul|\/ol)>/gi,'>\n<$1>')
    
    
    return html
}

// Clean html without replacing html completely and losing selection - used during editing
function cleanHTMLElements(el) {
    // Browsers tend to use style attributes to add all sorts of awful stuff to the html
    // No inline styles allowed
    el.find('p, div, b, i, h1, h2, h3, h4, h5, h6').removeAttr('style')
    el.find('span').removeAttr('style').removeAttr('lang')
    el.find('font').removeAttr('color')

    el.find()
}


function cleanAlign(cmd,sel) {
    console.log(sel.attr('style'))
    
    switch (cmd){
     case "justifyCenter": 
         
         if (sel.hasClass('align-center')) {
             sel.removeClass('align-center')
         } else {
             sel.addClass('align-center')
         }
         
         sel.removeClass('align-left').removeClass('align-right')
         sel.removeAttr('style') 
     break;
     case "justifyLeft": 
         if (sel.hasClass('align-left')) {
             sel.removeClass('align-left')
         } else {
             sel.addClass('align-left')
         }
     
         sel.removeClass('align-center').removeClass('align-right')
         sel.removeAttr('style') 
         
         
     break;
     case "justifyRight": 

         if (sel.hasClass('align-right')) {
             sel.removeClass('align-right')
         } else {
             sel.addClass('align-right')
         }

         sel.removeClass('align-center').removeClass('align-left')
         sel.removeAttr('style') 
     break;
     }
     
}

function startsWith(a,b) {
	if (a === undefined || b === undefined || a.length == 0 || b.length == 0) {
        return false;
    }
    result = a.lastIndexOf(b, 0)
    
    return (result.length > 0 && result == 0)
}



function getSelectionParentElement() {
    var parentEl = null, sel;
    if (window.getSelection) {
        sel = window.getSelection();
        if (sel.rangeCount) {
            parentEl = sel.getRangeAt(0).commonAncestorContainer;
            if (parentEl.nodeType != 1) {
                parentEl = parentEl.parentNode;
            }
        }
    } else if ( (sel = document.selection) && sel.type != "Control") {
        parentEl = sel.createRange().parentElement();
    }
    return parentEl;
}
