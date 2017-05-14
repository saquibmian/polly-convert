# polly-convert

polly-convert converts text to speech using Amazon Polly. An Amazon AWS account is required. 

Here is how to use polly-convert:

```
> polly-convert -file {filePath} -voice Amy -out {outputFile.mp3}
converting in 4 batch(es)
batch 1/4 complete
batch 2/4 complete
batch 3/4 complete
batch 4/4 complete
audio file written to output.mp3 in 6 seconds
```

## credentials

Credentials are pulled in from standard AWS sources. Run `aws configure` to set up your environment beforehand.
