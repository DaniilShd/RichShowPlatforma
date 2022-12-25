
    //Пример запроса на сервер для получения информации без перезагрузки страницы
    function updateLead() {
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
    updateLead()
    // setInterval(() => updateLead(), 1000);

    function updateStoreOrder() {
      fetch('http://localhost:8080/admin/fetch-orders')
    .then((response) => {
      return response.json();
    })
    .then((data) => {
      targetRaw = document.getElementById('new-order')
      targetConfirmed = document.getElementById('completed-order')
      targetArchive = document.getElementById('destroy-order')
      targetRaw.textContent = data['new-order']
      targetConfirmed.textContent = data['completed-order']
      targetArchive.textContent = data['destroy-order']
    });
  }
  updateStoreOrder()
  // setInterval(() => updateStoreOrder(), 1000);