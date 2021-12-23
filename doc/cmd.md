
commands:

- --version, V [Optional, Bool]: print the version
- --config, C [Must, String]: tool will first load the default flags from config file
- user flags [Optional]: tool will load user flags to overwrite the config
  - node, N [Single, String]: node name, which node to modify
  - --name [String]
  - --output_dir_path [String]
  - --output_remote [Bool]
  - --duration [Int]
  - --period [Int]
  - --infos [Multiple, String]
  - --diags [Multiple, String]