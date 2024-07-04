package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	orderedmap "github.com/elliotchance/orderedmap/v2"
	"github.com/stretchr/testify/assert"
	"kcl-lang.io/kpm/pkg/utils"
)

const testTomlDir = "test_data_toml"

func TestMarshalTOML(t *testing.T) {
	modfile := ModFile{
		Pkg: Package{
			Name:    "MyKcl",
			Edition: "v0.0.1",
			Version: "v0.0.1",
			Include: []string{"src/", "README.md", "LICENSE"},
			Exclude: []string{"target/", ".git/", "*.log"},
		},
		Dependencies: Dependencies{
			orderedmap.NewOrderedMap[string, Dependency](),
		},
	}

	dep := Dependency{
		Name:     "MyKcl1",
		FullName: "MyKcl1_v0.0.2",
		Source: Source{
			Git: &Git{
				Url: "https://github.com/test/MyKcl1.git",
				Tag: "v0.0.2",
			},
		},
	}

	ociDep := Dependency{
		Name:     "MyOciKcl1",
		FullName: "MyOciKcl1_0.0.1",
		Version:  "0.0.1",
		Source: Source{
			Oci: &Oci{
				Tag: "0.0.1",
			},
		},
	}

	modfile.Dependencies.Deps.Set("MyOciKcl1_0.0.1", ociDep)
	modfile.Dependencies.Deps.Set("MyKcl1_v0.0.2", dep)

	got_data := modfile.MarshalTOML()

	expected_data, _ := os.ReadFile(filepath.Join(getTestDir(testTomlDir), "expected.toml"))
	expected_toml := utils.RmNewline(string(expected_data))

	reversed_expected_data, _ := os.ReadFile(filepath.Join(getTestDir(testTomlDir), "expected_reversed.toml"))
	reversed_expected_toml := utils.RmNewline(string(reversed_expected_data))
	fmt.Printf("expected_toml: '%q'\n", expected_toml)
	fmt.Printf("reversed_expected_toml: '%q'\n", reversed_expected_toml)
	fmt.Printf("modfile: '%q'\n", got_data)
	fmt.Printf("expected_toml == got_data: '%t'\n", expected_toml == got_data)
	fmt.Printf("reversed_expected_toml == got_data: '%t'\n", reversed_expected_toml == got_data)
	assert.Equal(t, (utils.RmNewline(expected_toml) == utils.RmNewline(got_data)) ||
		(utils.RmNewline(reversed_expected_toml) == utils.RmNewline(got_data)), true, "got data is ", expected_toml, reversed_expected_toml, got_data)
}

func TestUnMarshalTOML(t *testing.T) {
	modfile := ModFile{}
	expected_data, _ := os.ReadFile(filepath.Join(getTestDir(testTomlDir), "expected.toml"))

	_ = toml.Unmarshal(expected_data, &modfile)
	fmt.Printf("modfile: %v\n", modfile)

	assert.Equal(t, modfile.Pkg.Name, "MyKcl")
	assert.Equal(t, modfile.Pkg.Edition, "v0.0.1")
	assert.Equal(t, modfile.Pkg.Version, "v0.0.1")
	assert.Equal(t, modfile.Pkg.Include, []string{"src/", "README.md", "LICENSE"})
	assert.Equal(t, modfile.Pkg.Exclude, []string{"target/", ".git/", "*.log"})
	assert.Equal(t, modfile.Dependencies.Deps.Len(), 2)
	assert.NotEqual(t, modfile.Dependencies.Deps.GetOrDefault("MyKcl1", TestPkgDependency), nil)
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Name, "MyKcl1")
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyKcl1", TestPkgDependency).FullName, "MyKcl1_v0.0.2")
	assert.NotEqual(t, modfile.Dependencies.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Source.Git, nil)
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Source.Git.Url, "https://github.com/test/MyKcl1.git")
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Source.Git.Tag, "v0.0.2")

	assert.NotEqual(t, modfile.Dependencies.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency), nil)
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Name, "MyOciKcl1")
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).FullName, "MyOciKcl1_0.0.1")
	assert.NotEqual(t, modfile.Dependencies.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Source.Registry, nil)
	assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Source.Registry.Tag, "0.0.1")
}

func TestMarshalLockToml(t *testing.T) {
	dep := Dependency{
		Name:     "MyKcl1",
		FullName: "MyKcl1_v0.0.2",
		Version:  "v0.0.2",
		Sum:      "hjkasdahjksdasdhjk",
		Source: Source{
			Git: &Git{
				Url: "https://github.com/test/MyKcl1.git",
				Tag: "v0.0.2",
			},
		},
	}

	ociDep := Dependency{
		Name:     "MyOciKcl1",
		FullName: "MyOciKcl1_0.0.1",
		Version:  "0.0.1",
		Sum:      "hjkasdahjksdasdhjk",
		Source: Source{
			Oci: &Oci{
				Reg:  "test_reg",
				Repo: "test_repo",
				Tag:  "0.0.1",
			},
		},
	}

	deps := Dependencies{
		orderedmap.NewOrderedMap[string, Dependency](),
	}

	deps.Deps.Set(dep.Name, dep)
	deps.Deps.Set(ociDep.Name, ociDep)
	tomlStr, _ := deps.MarshalLockTOML()
	expected_data, _ := os.ReadFile(filepath.Join(getTestDir(testTomlDir), "expected_lock.toml"))
	expected_toml := string(expected_data)
	assert.Equal(t, utils.RmNewline(expected_toml), utils.RmNewline(tomlStr))
}

func TestUnmarshalLockToml(t *testing.T) {
	deps := Dependencies{
		orderedmap.NewOrderedMap[string, Dependency](),
	}

	expected_data, _ := os.ReadFile(filepath.Join(getTestDir(testTomlDir), "expected_lock.toml"))
	expected_toml := string(expected_data)
	_ = deps.UnmarshalLockTOML(expected_toml)

	assert.Equal(t, deps.Deps.Len(), 2)
	assert.NotEqual(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency), nil)
	assert.Equal(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Name, "MyKcl1")
	assert.Equal(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).FullName, "MyKcl1_v0.0.2")
	assert.Equal(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Version, "v0.0.2")
	assert.Equal(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Sum, "hjkasdahjksdasdhjk")
	assert.NotEqual(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Source.Git, nil)
	assert.Equal(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Source.Git.Url, "https://github.com/test/MyKcl1.git")
	assert.Equal(t, deps.Deps.GetOrDefault("MyKcl1", TestPkgDependency).Source.Git.Tag, "v0.0.2")

	assert.NotEqual(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency), nil)
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Name, "MyOciKcl1")
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).FullName, "MyOciKcl1_0.0.1")
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Version, "0.0.1")
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Sum, "hjkasdahjksdasdhjk")
	assert.NotEqual(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Source.Oci, nil)
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Source.Oci.Reg, "test_reg")
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Source.Oci.Repo, "test_repo")
	assert.Equal(t, deps.Deps.GetOrDefault("MyOciKcl1", TestPkgDependency).Source.Oci.Tag, "0.0.1")
}

func TestUnMarshalTOMLWithProfile(t *testing.T) {
	modfile, err := LoadModFile(getTestDir("test_profile"))
	assert.Equal(t, err, nil)
	assert.Equal(t, modfile.Pkg.Name, "kpm")
	assert.Equal(t, modfile.Pkg.Version, "0.0.1")
	assert.Equal(t, modfile.Pkg.Edition, "0.0.1")
	assert.Equal(t, *modfile.Profiles.Entries, []string{"main.k", "xxx/xxx/dir", "test.yaml"})
}

func TestUnMarshalOciUrl(t *testing.T) {
	testDataDir := getTestDir("test_oci_url")

	testCases := []struct {
		Name          string
		DepName       string
		DepFullName   string
		DepVersion    string
		DepSourceReg  string
		DepSourceRepo string
		DepSourceTag  string
	}{
		{"unmarshal_0", "oci_pkg_name", "oci_pkg_name_0.0.1", "0.0.1", "ghcr.io", "test/helloworld", "0.0.1"},
		{"unmarshal_1", "oci_pkg_name", "oci_pkg_name_0.0.1", "0.0.1", "localhost:5001", "test/helloworld", "0.0.1"},
	}

	for _, tc := range testCases {
		modfile, err := LoadModFile(filepath.Join(testDataDir, tc.Name))
		assert.Equal(t, err, nil)
		assert.Equal(t, modfile.Dependencies.Deps.Len(), 1)
		assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("oci_pkg_name", TestPkgDependency).Name, tc.DepName)
		assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("oci_pkg_name", TestPkgDependency).FullName, tc.DepFullName)
		assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("oci_pkg_name", TestPkgDependency).Version, tc.DepVersion)
		assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("oci_pkg_name", TestPkgDependency).Source.Oci.Reg, tc.DepSourceReg)
		assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("oci_pkg_name", TestPkgDependency).Source.Oci.Repo, tc.DepSourceRepo)
		assert.Equal(t, modfile.Dependencies.Deps.GetOrDefault("oci_pkg_name", TestPkgDependency).Source.Oci.Tag, tc.DepVersion)
	}
}

func TestMarshalOciUrl(t *testing.T) {
	testDataDir := getTestDir("test_oci_url")

	expectPkgPath := filepath.Join(testDataDir, "marshal_0", "kcl_mod_bk")
	gotPkgPath := filepath.Join(testDataDir, "marshal_0", "kcl_mod_tmp")

	expect, err := LoadModFile(expectPkgPath)
	assert.Equal(t, err, nil)

	err = os.MkdirAll(gotPkgPath, 0755)
	assert.Equal(t, err, nil)
	gotFile, _ := os.Create(filepath.Join(gotPkgPath, "kcl.mod"))

	defer func() {
		err = gotFile.Close()
		assert.Equal(t, err, nil)
		err = os.RemoveAll(gotPkgPath)
		assert.Equal(t, err, nil)
	}()

	modfile := ModFile{
		Pkg: Package{
			Name:    "marshal_0",
			Edition: "v0.9.0",
			Version: "0.0.1",
		},
		Dependencies: Dependencies{
			orderedmap.NewOrderedMap[string, Dependency](),
		},
	}

	ociDep := Dependency{
		Name:     "oci_pkg",
		FullName: "oci_pkg_0.0.1",
		Version:  "0.0.1",
		Source: Source{
			Oci: &Oci{
				Reg:  "ghcr.io",
				Repo: "kcl-lang/oci_pkg",
				Tag:  "0.0.1",
			},
		},
	}

	modfile.Dependencies.Deps.Set("oci_pkg_0.0.1", ociDep)

	got_data := modfile.MarshalTOML()
	_, err = gotFile.WriteString(got_data)
	assert.Equal(t, err, nil)

	got := ModFile{}
	err = got.LoadModFile(filepath.Join(gotPkgPath, "kcl.mod"))
	assert.Equal(t, err, nil)

	assert.Equal(t, expect.Pkg.Name, got.Pkg.Name)
	assert.Equal(t, expect.Pkg.Edition, got.Pkg.Edition)
	assert.Equal(t, expect.Pkg.Version, got.Pkg.Version)
	assert.Equal(t, expect.Dependencies.Deps.Len(), got.Dependencies.Deps.Len())
	assert.Equal(t, expect.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Name, got.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Name)
	assert.Equal(t, expect.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).FullName, got.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).FullName)
	assert.Equal(t, expect.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.Oci.Reg, got.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.Oci.Reg)
	assert.Equal(t, expect.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.Oci.Repo, got.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.Oci.Repo)
	assert.Equal(t, expect.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.Oci.Tag, got.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.Oci.Tag)
	assert.Equal(t, expect.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.IntoOciUrl(), got.Dependencies.Deps.GetOrDefault("oci_pkg", TestPkgDependency).Source.IntoOciUrl())
}

func TestMarshalOciUrlIntoFile(t *testing.T) {
	testDataDir := getTestDir("test_oci_url")

	testCases := []string{"marshal_2"}

	for _, tc := range testCases {
		readKclModPath := filepath.Join(testDataDir, tc)
		modfilePath := filepath.Join(readKclModPath, "kcl.mod")
		expectPath := filepath.Join(readKclModPath, "expect.mod")

		readKclModFile := ModFile{}
		err := readKclModFile.LoadModFile(modfilePath)
		assert.Equal(t, err, nil)
		writeKclModFileContents := readKclModFile.MarshalTOML()
		expectKclModFileContents, err := os.ReadFile(expectPath)
		assert.Equal(t, err, nil)

		assert.Equal(t, utils.RmNewline(string(expectKclModFileContents)), utils.RmNewline(writeKclModFileContents))
	}
}
