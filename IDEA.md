It is common in every language that some words occur more often than others. For natural learners who learn by exposing themself to the language rather than by memorizing it, it is natural that they remember more frequent words faster. Such natural approach is called language acquisition and is based on learning by listening or reading rather than by learning grammar and memorizing vocabulary. Moreover, some polyglots say that they need to hear the word at least 20 to 30 times to actually remember it. We will use this in building our theory.

Let's say there are words $S_{1}$..$S_{N}$ in a language with probabilities of occurence $p_{1}$..$p_{N}$. 
Let's say a learner needs to see a word $k$ times to remember it.

Denote $X_{in}$ the number of occurencies of the word $S_{i}$ in a random sample of a language text consisting of n words.

Let's assume that some sample person is familiar only with this random language text sample. Then the probability that this person remembers/knows word $S_{i}$ is 
```math
$$P(X_{in} >= k) = 1 - P(X_{in} < k)$$, where $$X_{in}$$ has binomial distribution.
```

Let's create sample of J people which are familiar with samples of text of size $n_{1} .. n_{J}$ and 
generate their knowledge of all words based on calculated probabilities. Such people have different total number of known words and it holds that they know more frequent words with higher probability. We predict how many words a person knows based on comparison with this sample of generated people.

