
# This module is for creating go packages

    export def main [] { bat scripts.nu }


    # Creates a normal package
    export def "main create" [package_name:string] {

        mkdir $package_name

        $"package ($package_name)\n\n\n func main\(){\n\n\n}"
        | save $"($package_name)/main.go"

    }

    # Deletes a package
    export def "main delete" [package_name:string] {

        rm $package_name -r

    }

    # Creates a package using `bubbletea-model-template.txt`.
    # To scaffold the file
    export def "main create-btf" [package_name:string] {

        open bubbletea-model-template.txt
        | collect
        | str replace "main" $package_name
        | save $"($package_name)/main.go"

    }


