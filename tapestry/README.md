# Tapestry (TA implementation)

## Note
By using this implementation, you have signed the [code exchange contract](https://forms.gle/CMdu42Wpf5cMos6n8)
and agree not to share it with anyone else.

## Instructions for Usage in Puddlestore

You will need to use your own Tapestry implementation or the TA implementation of Tapestry for this project (see the handout for more details). However, you will not be using import statements in Go for this. Instead, you will be using a powerful feature known as [the `replace directive` in Go modules](https://thewebivore.com/using-replace-in-go-mod-to-point-to-your-local-module/)

### Steps:

1. Download the Tapestry implementation to a local folder. 
2. Update this line in `go.mod`

```
replace tapestry => /path/to/your/tapestry/implementation/root/folder
```

so that imports of `tapestry` now point to your local folder where you've downloaded Tapestry. 

That's it! When running tests in Gradescope, we will automatically rewrite this line to point to our TA implementation, so you will be tested against our implementation and not be penalized for any issues of your own. 



## Feedback
Your feedback are always welcome! If you find any inconsistencies or bugs in the TA implementation,
please post on Edstem or fill out [the anonymous feedback form](http://cs.brown.edu/courses/cs138/s22/feedback.html)