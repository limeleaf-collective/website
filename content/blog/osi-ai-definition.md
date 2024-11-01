+++
title = "The OSI defined open-source AI. Surprise, big-tech AI projects aren't open source."
date = 2024-10-31
draft = false
authors = ["John Luther"]

[taxonomies]
tags = ["Company", "AI"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

This week, the [Open Source Initiative](https://opensource.org/) (OSI) released ["The Open Source AI Definition"](https://www.prweb.com/releases/the-open-source-initiative-announces-the-release-of-the-industrys-first-open-source-ai-definition-302288553.html) (OSAID). This is an important step in establishing guidelines around what "AI" means in relation to intellectual property rights and how open-source definitions apply to the technology.  

<!-- more -->

When it comes to standards bodies, the OSI is about as legit as they come. [In 1998](https://opensource.org/about), they [literally](https://www.rd.com/article/what-does-literally-mean/) coined the term *open source* and have since served as the stewards of [The Open Source Definition](https://opensource.org/osd).

Unsurprisingly, the self-proclaimed "open source" AI systems that are funded by tech megacorporations and VCs don't comply with the OSAID. To quote the [OSI's FAQ](https://hackmd.io/@opensourceinitiative/osaid-faq): 

> "Those that have been analyzed and don't pass because they lack required components and/or their legal agreements are incompatible with the Open Source principles: Llama2 (Meta), Grok (X/Twitter), Phi-2 (Microsoft), Mixtral (Mistral)." 
 
(Note that OpenAI's models aren't relevant here because, despite their name, they don't claim to be open source.)

The FAQ doesn't give specifics for why any of the systems don't comply, but I imagine the data disclosure requirements will be a major blocker for them. To comply with the OSAID, an AI system must disclose "sufficiently detailed information about the data used to train the system." They define this as:

> "(1) the complete description of all data used for training, including (if used) of unshareable data, disclosing the provenance of the data, its scope and characteristics, how the data was obtained and selected, the labeling procedures, and data processing and filtering methodologies; (2) a listing of all publicly available training data and where to obtain it; and (3) a listing of all training data obtainable from third parties and where to obtain it, including for fee."

As we have seen, AI companies [have played fast and loose](https://www.theverge.com/23444685/generative-ai-copyright-infringement-legal-fair-use-training-data) [with copyrighted content](https://www.proofnews.org/apple-nvidia-anthropic-used-thousands-of-swiped-youtube-videos-to-train-ai/) and [user data](https://telehealth.org/facebooks-ai-training-controversy-the-ethical-implications-of-using-childrens-photos/) to train their models. I could be wrong, but I don't expect they will become OSAID-compliant ([or endorse the OSAID](https://opensource.org/ai/endorsements)) any time soon.
