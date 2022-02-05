# API Conversão monetária

API para realizar trocas entre moedas  

### Rotas
```
GET /   retorna a base, data e taxa das moedas disponiveis para troca
GET /moedas   retorna as moedas disponiveis para troca
GET /moedas/{moeda}   retorna o valor da moeda passada
GET /moedas/troca   retorna uma string com a conversão entre as moedas passadas

GET /usuarios   retorna todos os usuários cadastrados e suas trocas
GET /usuario/{id}   retorna as trocas do usuário passado
POST /usuario/troca   recebe as moedas e quantidade para a troca, salva um objeto de troca no usuário e retorna o objeto da troca (necessario passar 'bearer token')
DELETE /usuario/troca   exclui uma troca com base na data passada e retorna as trocas do usuário (necessario passar 'bearer token')

POST /login   recebe um 'username' e 'password' e retorna um JWT 
POST /registrar   recebe um 'username' e 'password', cria um usuário e retorna um JWT 
DELETE /excluir   exclui um usuário 
```

<br>

#### extras
* Banco de dados utilizado - [MongoDB Atlas](https://www.mongodb.com/atlas/database)
* valores monetários pegos de [https://www.currency-api.com](https://www.currency-api.com)
