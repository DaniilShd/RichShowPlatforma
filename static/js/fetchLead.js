
    //Пример запроса на сервер для получения информации без перезагрузки страницы
    function sayHi() {
        fetch('http://localhost:8080/admin/fetch-leads')
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        targetRaw = document.getElementById('raw-leads')
        targetConfirmed = document.getElementById('confirmed-leads')
        targetArchive = document.getElementById('archive-leads')
        targetRaw.textContent = data['raw-leads']
        targetConfirmed.textContent = data['confirmed-leads']
        targetArchive.textContent = data['archive-leads']
      });
    }
    sayHi()
    setInterval(() => sayHi(), 1000);