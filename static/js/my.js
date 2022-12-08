function add_field(){

    var x = document.getElementById("many_input");

    var new_div = document.createElement("div");
    new_div.setAttribute("class", "row");
    // создаем новое поле ввода
    var new_field = document.createElement("input");
    // установим для поля ввода тип данных 'text'
    new_field.setAttribute("type", "text");
    // установим имя для поля ввода
    new_field.setAttribute("name", "name_of_points[]");
    new_field.setAttribute("class", "form-control field");
  
    new_field.setAttribute("id", "name_of_points");
    // определим место вствки нового поля ввода (перед каким элементом его вставить)
    var pos = x.childElementCount;
    new_div.appendChild(new_field)
    // добавим поле ввода в форму
    // x.insertBefore(new_div, x.childNodes[pos]);
    x.appendChild(new_div)
    }

    function delete_field(){
    let x = document.getElementsByClassName('field')
    x[x.length-1].remove()     
    }