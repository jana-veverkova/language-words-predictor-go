# Language words predictor

Language word predictor consists of several main packages:

1) **frequencydictionary** - enables to create a frequency dictionary from any txt files saved in data/original/'target language' containing any text in the target language. Frequency dictionary is then stored as data/frequencyDictionaries/'target language'.csv. This file contains words and number of their occurence in the source files ordered from the most common one.

Following rules to the words in the source files are implied:
- Special characters and numbers at the beginning or the end of the word are trimmed.
- When a special character or a number occurs in the middle of the word the whole word is omitted.
- Characters ' and - in the middle of words are allowed. I.e. words like "it's" or "sour-milk" will appear in the frequency dictionary.
- The words in the frequency dictionary will be lower cased.

2) **populationsample** - generates samples of people while for every person and each of the first 30000 words value 0 or 1 is saved. 0 means that person doesn't know the word, 1 mean that person knows the word. Population samples are then stored to data/processed/populationSamples directory. This sampling is based on theory explained in IDEA.md.

3) **traintestsplit** - splits population sample randomly 1:4 into a train and a test set. The train and the test sets are saved to data/processed diresctory.

4) **model**
- enables to run prediction of number of words a user knows in the terminal. User will be asked several words (max 120) based on frequency dictionary and at the end the prediction is displayed.
- enables to train the setting parameters of the prediction model. For given set of parameters for every person from the test population sample the prediction is run and the average accuracy and mistake size is displayed.
- the prediction algorithm is independent of a frequency dictionary file and can be used with any frequency dictionary.


