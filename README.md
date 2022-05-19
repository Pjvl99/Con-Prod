# productor-consumidor proyecto

- Jean Pierre Mejicanos
- Adriana Mundo
- Pablo Velasquez

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


En la base de datos se pueden obervar los resultados, sin que haya alguna actriz repetida 

![img3](https://user-images.githubusercontent.com/61527863/169353084-1d0c4234-7444-4780-975e-5e7174d47c93.jpeg)

