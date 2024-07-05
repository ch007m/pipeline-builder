## Pipeline builder POC

To build the executable
```bash
make
```

Next, launch it
```bash
./out/pipe-builder -h
A tekton pipeline builder able to create from templates and a configurator files pipelines ans tasks.

Usage:
  pipeline-builder [flags]

Flags:
  -c, --configurator string   path of the configurator file
  -h, --help                  help for pipeline-builder
  -o, --output string         Output where pipelines should be saved
```  
  
If there is a configuration file `conf.yaml` created at the root of this project and that you want to generate the pipelines yaml files under `out/flows`, then execute this command:
```bash
./out/pipe-builder -c conf.yaml -o out/flows
```
Next, check the pipeline(s) generated under `./out/flows`

