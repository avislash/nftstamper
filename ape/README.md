# APE Bot
[![Twitter Follow](https://img.shields.io/twitter/follow/avis1ash?style=social)](https://twitter.com/avis1ash "Follow me on Twitter!")
[![Discord](https://img.shields.io/badge/avislash%235874-7289DA?logo=discord&logoColor=white)](#) 
[![ENS](https://img.shields.io/badge/ENS-avislash.eth-blueviolet?logo=ethereum)](https://avislash.eth.xyz/)
  

The APE Bot is an instantiation of the NFT Stamper designed to support Applied Primate Engineering/Fragment collections
<div align="center">
  <img src="https://github.com/avislash/nftstamper/blob/main/example_images/sentinels_example.jpg" alt="Example Image" width="400"/>
</div>

# Installation
Standard go installation
```
$ go get github.com/avislash/nftstamper@latest
```

# 
The Bot Configuration is driven through the [config.yaml file](https://github.com/avislash/nftstamper/blob/main/ape/config.yaml). 

The `ipfs_endpoint` is the endpoint to an IPFS node that can be used to retrive files off IPFS. The endpoint format must be in [multiaddr format](https://github.com/multiformats/multiaddr#encoding).
 - Note: While any available IPFS endpoint can be specified for better performance consider hosting and running a local IPFS node. Refer to the instructions found in the [Kubo IPFS repository](https://github.com/ipfs/kubo) for how to download, install, and configure a local IPFS node.

The `log_level` option is used to specify the logging level. Valid options are `debug`, `info`, and `error`. All logging defaults to info unless specified otherwise in the config file.

The `metadata_endpoint` is the primary web endpoint for scraping metadata against.

The `discord_bot_token` is the bot API key to use. It can be either entered in the config file or saved to the environment aswhatever value is specified for this entry

The `image_procesor` section is for configuring the image processor and defining mappings of characteristics to overlay images. 

## Image Processor Mappings
The `image_processor_mappings` part of the `image_procesor` section is used to define mappings for the various `gm` command. 


### GM Mappings
The `gm_mappings` part specifically maps base armors to mug hands. The keys under `gm_mappings` are the names of different base armors (like `trippy` in the example), and the values are paths to the corresponding images (lab for lab mugs and path for path mugs).

```yaml
image_processor_mappings:
    gm_mappings:
            trippy: 
               lab: ./ape/gm_assets/TRIPPY.png
               path: ./ape/gm_assets/TRIPPY_PATH.png
```
In this example, for the base armor named trippy, there are two images: TRIPPY.png which is used for the lab, and TRIPPY_PATH.png which is used for the path.

### Filters
The `filters` part of the image_procesor section is used to define default opacity levels for the GM smoke, as well as trait-specific weights that can be applied to the opacity. The opacity key is used specifically for the gm_smoke and accepts valid values between 0 and 1, where 1 is 100% opacity and 0.5 is 50% opacity, and so forth.

```yaml
filters:
    opacity:
        trippy: 
           name: trippy
           Default: 0.9
           Weights:
               path robe: 1
```
Here, the Default key under trippy sets the default opacity for the trippy base armor to 0.9 (or 90%). The Weights key is used to define trait-specific opacity levels. For instance, the path robe trait will have an opacity of 1 (or 100%) when applied.

# Usage & Examples
Before running the bot ensure that the `APE_DISCORD_BOT` environemnt variable is set.

To run the bot use
```
./nftstamper apebot
```

This will instantiate the NFT Stamper to use mappings for Applied Primate Engineering for its image processing and other commands.


## Supported Collections

- [Applied Primate Engineering](https://github.com/avislash/nftstamper/tree/main/ape)
   - MegaForce Sentinels [![MegaForce Sentinels](https://img.shields.io/badge/Supported-100%25-brightgreen)](#)


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
 - LetMeDo
    - Thank you for inspiring this entire project and bot.
 - [The Best Online README Generator](https://readme.so/)

## Contact Me <a name="contact-me"></a>
[![twitter](https://img.shields.io/badge/@avis1ash-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white)](https://twitter.com/avis1ash)
[![Discord](https://img.shields.io/badge/avislash%235874-7289DA?style=for-the-badge&logo=discord&logoColor=white)](#)
[![ENS](https://img.shields.io/badge/ENS-avislash.eth-blueviolet?style=for-the-badge&logo=ethereum)](https://avislash.eth.xyz/)
