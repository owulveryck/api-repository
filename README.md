Still a work in progress...

See the `doc` directory for a work in progress of the slides
to run it:
```
# go get github.com/owulveryck/api-repository
# present
```

Note: The slides can be viewed online thanks to the godoc project: 
[here](https://talks.godoc.org/github.com/owulveryck/api-repository/doc/bbl.slide)

# Abstract

Web Oriented API is a must for a company that wants to interact with the outside world.
Most of the companies tend to develop an API as an entry point to their platform. On top of that, data is in the center of the strategies; therefore grabbing data through API is also fundamental.
The *availability* of the API is more than ever a concept that can cause severe business loss if not handled properly.

As an SRE, part of a product team, I expose what I expect from an API in term of objectives to manage and pilot the reliability of the application.

Then I use those objectives to design and develop an underlying API that can handle a bulk of product to store them in a DAO.

In the end, I give some pointers to host this "_cloud-native_" application on the cloud in a cost-effective way (using GCP app engine).


## Takeaway

Unit Tests are fundamental for developing software in quality. Using objectives is a must to think about reliability, and it starts from the design of the application.

# Résumé

Les API sont devenues un obligation pour les sociétés qui souhaitent s'ouvrir à l'extérieur.

Ainsi, c'est pour beaucoup une porte d'entrée vers les plateformes de services. De plus, la donnée est au centre des stratégie. C'est donc naturellement que l'API s'est imposée comme un moyen de récupérer de la donnée dans le but d'y ajouter de la valeur.

La *disponibilité* de l'API est donc, plus que jamais, un concept qui revet un importance particulière si elle n'est pas prise au sérieux. Elle peut être la cause de perte de business importante.

En tant que SRE dans une équipe produit, j'expose dans ce talk ce que j'attend d'une API en terme d'objectis. Le but est de piloter la fiabilité de l'application.

Ensuite j'utilise ces objectifs pour designer et développer un moteur d'API (en Go) capable de prendre en charge un bulk de produits et de le stocker dans un DAO.

Enfin, je donne quelques pointers pour héberger cette application "cloud-native" de manière efficace et peu onéreuse sur le cloud public GCP.

## Takeaway

Les tests unitaires sont fondamentaux pour developper un logiciel de qualité. Utiliser des objectifs" est un must pour penser à la **fiabilité** de l'application et ceci démarre dès la phase de conception.
