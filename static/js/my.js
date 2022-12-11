function add_field(){

    var x = document.getElementById("many_input");

    var new_div = document.createElement("div");
    new_div.setAttribute("class", "field");
    // создаем новое поле ввода
    var new_field = document.createElement("input");
    // установим для поля ввода тип данных 'text'
    new_field.setAttribute("type", "text");
    // установим имя для поля ввода
    new_field.setAttribute("name", "name_of_points[]");
    new_field.setAttribute("class", "form-control");
  
    new_field.setAttribute("id", "name_of_points");
    // определим место вствки нового поля ввода (перед каким элементом его вставить)
    // var pos = x.childElementCount;
    new_div.appendChild(new_field)
    // добавим поле ввода в форму
    // x.insertBefore(new_div, x.childNodes[pos]);
    x.appendChild(new_div)
    }

    function delete_field(){
    let x = document.getElementsByClassName('field')
    x[x.length-1].remove()     
    }


    function add_item(){
        
        var child = document.getElementById("item_many").cloneNode(true);
        child.style.display = "inline";
        var main = document.getElementById("item-main");
        main.appendChild(child)
        }
    
        function delete_item(){
        let x = document.getElementsByClassName('item-store')
        if (x.length == 1) {
            return
        }
        x[x.length-1].remove()     
        }


//Lead start-----------------------------------------------------------------------------------------------------------------------------
function add_program(){

    var child = document.getElementById("many_program").cloneNode(true);
    child.style.display = "inline";
    var main = document.getElementById("item-main-program");
    main.appendChild(child)
    }

    function delete_program(){
    let x = document.getElementsByClassName('item-program')
    if (x.length == 1) {
        return
    }
    x[x.length-1].remove()     
    }


function add_master_class(){

    var child = document.getElementById("many_master_class").cloneNode(true);
    child.style.display = "inline";
    var main = document.getElementById("item-main-master_class");
    main.appendChild(child)
    }

    function delete_master_class(){
    let x = document.getElementsByClassName('item-master_class')
    if (x.length == 1) {
        return
    }
    x[x.length-1].remove()     
    }


function add_animation(){

var child = document.getElementById("many_animation").cloneNode(true);
child.style.display = "inline";
var main = document.getElementById("item-main-animation");
main.appendChild(child)
}

function delete_animation(){
let x = document.getElementsByClassName('item-animation')
if (x.length == 1) {
    return
}
x[x.length-1].remove()     
}


function add_party_and_quest(){

    var child = document.getElementById("many_party_and_quest").cloneNode(true);
    child.style.display = "inline";
    var main = document.getElementById("item-main-party_and_quest");
    main.appendChild(child)
    }
    
    function delete_party_and_quest(){
    let x = document.getElementsByClassName('item-party_and_quest')
    if (x.length == 1) {
        return
    }
    x[x.length-1].remove()     
    }


function add_other(){

var child = document.getElementById("many_other").cloneNode(true);
child.style.display = "inline";
var main = document.getElementById("item-main-other");
main.appendChild(child)
}

function delete_other(){
let x = document.getElementsByClassName('item-other')
if (x.length == 1) {
    return
}
x[x.length-1].remove()     
}

function add_hero(){

    var child = document.getElementById("many_hero").cloneNode(true);
    child.style.display = "inline";
    var main = document.getElementById("item-main-hero");
    main.appendChild(child)
}
    
function delete_hero(){
    let x = document.getElementsByClassName('item-hero')
    if (x.length == 1) {
        return
    }
    x[x.length-1].remove()     
}


function add_assistant(){

    var child = document.getElementById("many_assistant").cloneNode(true);
    child.style.display = "inline";
    var main = document.getElementById("item-main-assistant");
    main.appendChild(child)
}
    
function delete_assistant(){
    let x = document.getElementsByClassName('item-assistant')
    if (x.length == 1) {
        return
    }
    x[x.length-1].remove()     
}
//Lead end-----------------------------------------------------------------------------------------------------------------------------