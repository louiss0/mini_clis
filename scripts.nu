
# This module is for creating go packages

def validate_path_exists [path: string] {
    if ($path | path exists) {
        error make {msg: $"($path) already exists"}
    }
}

def validate_package_name [name: string] {
    # Create regex pattern for valid Go package names
    let pattern = '^[a-z][a-z0-9_]*$'
    
    # List of Go keywords that can't be used
    let keywords = [
        'break' 'default' 'func' 'interface' 'select'
        'case' 'defer' 'go' 'map' 'struct'
        'chan' 'else' 'goto' 'package' 'switch'
        'const' 'fallthrough' 'if' 'range' 'type'
        'continue' 'for' 'import' 'return' 'var'
    ]

    # Check if it's a keyword
    if ($keywords | any { |k| $k == $name }) {
        error make {msg: $"($name) is a Go keyword and cannot be used as a package name"}
    }

    # Validate against pattern
    if not ($name | str replace -r $pattern "") == "" {
        error make {msg: $"Invalid package name: ($name). Package names must start with a letter and contain only lowercase letters, numbers, or underscores"}
    }
}

export def main [] { bat scripts.nu }


    # Creates a normal package
    export def "main create" [package_name:string] {
        validate_package_name $package_name
        validate_path_exists $package_name
        mkdir $package_name
        $"package ($package_name)\n\nfunc main() {\n\n}\n"
        | save $"($package_name)/main.go"

    }

    # Deletes a package
    export def "main delete" [package_name:string] {

        rm $package_name -r

    }

    # Creates a package using `bubbletea-model-template.txt`.
    # To scaffold the file
    export def "main create-btf-package" [package_name:string] {
        validate_package_name $package_name
        validate_path_exists $package_name
        let go_main_function_snippet = "
    func main() {

	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {

		log.Fatal(err)

	}

}
"

        open bubbletea-model-template.txt
        | collect
        | str replace "main" $package_name
        | $"($in)\n($go_main_function_snippet)"
        | save $"($package_name)/main.go"

    }

    export def "main create-btf-file" [filename:string package_name:string] {
        validate_package_name $package_name
        validate_path_exists $"($package_name)/($filename).go"
        open bubbletea-model-template.txt
        | collect
        | str replace "main" $package_name
        | save $"($package_name)/($filename).go"


    }
