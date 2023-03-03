// simple client implementation
const sc = new EventSource('http://localhost:8080/events');

sc.onmessage = (e) => {
  const payload = JSON.parse(e.data);

  document.getElementById('sent').innerText = payload.sent;
  document.getElementById('received').innerText = payload.received;
  document.getElementById('total').innerText = payload.total;
};
