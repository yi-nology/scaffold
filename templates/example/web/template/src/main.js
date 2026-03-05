import './style.css'

document.querySelector('#app').innerHTML = `
  <div class="container">
    <h1>{{.project_name}}</h1>
    <p>{{.description}}</p>
  </div>
`
