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