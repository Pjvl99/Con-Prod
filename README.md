## Web scrapping utilizando productores y consumidores


**Descripción:** 
Se requiere almacenar en una base de datos la relación existente entre películas y actrices.
Los films en los que ha laborado cada actriz están definidos en cada página web de la filmografía en la que la actriz ha actuado. Se debe de leer el HTML del sitio de actrices y almacenar lo que se le indica  en una base de datos mysql 




### Tecnologías utilizadas:
- Go: para crear y manejar los consumidores y productores
- Docker
    - mysql: para guardar todos los datos   

## Ejecución del proyecto
#### Base de datos - Docker Mysql
```
docker-compose up
``` 

#### GO - productores y consumidores
```
./pcb #Productores #Consumidores
```
En #Productores se pone la cantidad de productores y en #Consumidores la cantidad de consumidores. Para ejecutar esto el docker debe estar corriendo.

### Ejecución
La idea es crear una cantidad de productores que recojan los datos de cada actriz y pasarsela a los consumidores sin que ocurra una condición de carrera que haga que los datos se repitan.

![img1](https://user-images.githubusercontent.com/61527863/169353015-8e18b421-fc4e-4f4c-aeff-4f6e19153791.jpeg)

![img4](https://user-images.githubusercontent.com/61527863/169354683-68796bdd-935c-4026-af43-6893058016b9.jpeg)

![img5](https://user-images.githubusercontent.com/61527863/169354693-b7e2aa4a-4b3b-4ca9-98ae-b7e345f41d8c.jpeg)

![img6](https://user-images.githubusercontent.com/61527863/169354743-bef545b1-45b0-4497-9e5d-04aa796a8cac.jpeg)

En la base de datos se pueden obervar los resultados, sin que haya alguna actriz repetida 

![img2](https://user-images.githubusercontent.com/61527863/169353632-0404ee0a-17f4-4e96-8b68-fe7c730d2426.jpeg)

![img3](https://user-images.githubusercontent.com/61527863/169354305-68f2977e-48ef-437b-80a4-234268f83136.jpeg)

