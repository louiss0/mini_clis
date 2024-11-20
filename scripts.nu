
# This module is for creating go packages
module package {


    export def main [] { help package }


    # Creates a mormal package
    export def create [package_name:string] {

        mkdir $package_name

        $"package ($package_name)\n\n\n func main\(){\n\n\n}"
        | save $"($package_name)/main.go"

    }

    # Deletes a package
    export def delete [package_name:string] {

        rm $package_name -r

    }

    # Creates a package using `bubbletea-model-template.txt`.
    # To scaffold the file
    export def create-btf [package_name:string] {

        open bubbletea-model-template.txt
        | collect
        | str replace "main" $package_name
        | save $"($package_name)/main.go"

    }

}

overlay use --prefix package
