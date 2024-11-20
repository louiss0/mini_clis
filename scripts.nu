
module package {


 export def main [] { help package }


 export def create [package_name:string] {

    mkdir $package_name

    $"package ($package_name)\n\n\n func main\(){\n\n\n}"
    | save $"($package_name)/main.go"

}

export def delete [package_name:string] {

 rm $package_name -r

}

}

overlay use --prefix package
