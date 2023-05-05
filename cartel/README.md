
# Cartel Bot
[![Twitter Follow](https://img.shields.io/twitter/follow/avis1ash?style=social)](https://twitter.com/avis1ash "Follow me on Twitter!")
[![Discord](https://img.shields.io/badge/avislash%235874-7289DA?logo=discord&logoColor=white)](#) 
[![ENS](https://img.shields.io/badge/ENS-avislash.eth-blueviolet?logo=ethereum)](https://avislash.eth.xyz/)
  

The Cartel Bot is an instantiation of the NFT Stamper designed to support Mutant Cartel/Novel Lab collections
<div align="center">
  <img src="https://github.com/avislash/nftstamper/blob/feature/README/example_images/example.jpg" alt="Example Image" width="400"/>
</div>

# Installation
Standard go installation
```
$ go get github.com/avislash/nftstamper@latest
```
The bot currently assumes acess to a locally running IPFS node for any IPFS retrievals it may need to perform.

Refer to the instructions found in the [Kubo IPFS repository](https://github.com/ipfs/kubo) for how to download, install, and configure a local IPFS node.


# Configuration
The Bot Configuration is driven through the [config.yaml file](https://github.com/avislash/nftstamper/blob/main/cartel/config.yaml). 

The `metadata_endpoint` is the primrary web endpoint for scraping metadata against.

The `image_procesor` section is for configuring the image processor and defining mappings of characteristics to overlay images. 
 - The `gm_mappings` section is used to map Hound Background Traits to [overlay images](https://github.com/avislash/nftstamper/tree/main/cartel/bowls). These mappings directly impact the output of the `/gm` command

# Usage & Examples
Before running the bot ensure that the `CARTEL_DISCORD_BOT` environemnt variable is set.

To run the bot use
```
./nftstamper cartelbot
```

This will instantiate the NFT Stamper to use mappings for The Mutant Cartel for its image processing and other commands.


## Supported Collections

- [The Mutant Cartel](https://github.com/avislash/nftstamper/tree/main/cartel)
   - Mutant Hounds [![Mutant Hounds](https://img.shields.io/badge/Supported-90%25-yellow)](#)
     - Mega Hound Support Pending


## How do I get a collection integrated
You may add the collection yourself by either
1. Cloning the repo and modifying the source code and/or artwork 
2. Contacting me using one of the options listed in the [Contact Me](#contact-me) section 
3. Opening an issue via the [Github Issue Tracker](https://github.com/avislash/nftstamper/issue)


## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Acknowledgements
 - [Line Hammett](https://linehammett.com/about-me)
    - Thank you for the amazing Mutant Hound bowl artwork you designed and contributed. Please checkout [her otherwork here](https://linehammett.com/digiart-and-nfts)
 - [The Best Online README Generator](https://readme.so/)

## Contact Me <a name="contact-me"></a>
[![twitter](https://img.shields.io/badge/@avis1ash-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/avis1ash)
[![Discord](https://img.shields.io/badge/avislash%235874-7289DA?style=for-the-badge&logo=discord&logoColor=white)](#)
[![ENS](https://img.shields.io/badge/ENS-avislash.eth-blueviolet?style=for-the-badge&logo=ethereum)](https://avislash.eth.xyz/)

## TODO
1. Integrate Vyper for better config/env parsing
2. Allow for setting IPFS Endpoint via config/env variable
3. Scrape MH Metadata endpoint per token from The Mutant Hounds smart contract
4. Update IPFS client to parse Files/Metadata + Images
5. Update Metadata Fetcher to fetch from either IPFS and/or Web Endpoint
5. Consider refactoring and abstracting both to lib