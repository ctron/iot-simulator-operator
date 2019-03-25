package e2e

import (
	"testing"

	"github.com/coreos/go-semver/semver"
	"github.com/stretchr/testify/require"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"
	"github.com/operator-framework/operator-lifecycle-manager/pkg/controller/registry"
	packagev1alpha1 "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/packagemanifest/v1alpha1"
	pmversioned "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/client/clientset/versioned"
)

type packageManifestCheckFunc func(*packagev1alpha1.PackageManifest) bool

func packageManifestHasStatus(pm *packagev1alpha1.PackageManifest) bool {
	// as long as it has a package name we consider the status non-empty
	if pm == nil || pm.Status.PackageName == "" {
		return false
	}

	return true
}

func fetchPackageManifest(t *testing.T, pmc pmversioned.Interface, namespace, name string, check packageManifestCheckFunc) (*packagev1alpha1.PackageManifest, error) {
	var fetched *packagev1alpha1.PackageManifest
	var err error

	err = wait.Poll(pollInterval, pollDuration, func() (bool, error) {
		t.Logf("Polling...")
		fetched, err = pmc.Packagemanifest().PackageManifests(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		return check(fetched), nil
	})

	return fetched, err
}

func TestPackageManifestLoading(t *testing.T) {
	// create a simple catalogsource
	packageName := genName("nginx")
	stableChannel := "stable"
	packageStable := packageName + "-stable"
	manifests := []registry.PackageManifest{
		registry.PackageManifest{
			PackageName: packageName,
			Channels: []registry.PackageChannel{
				registry.PackageChannel{Name: stableChannel, CurrentCSVName: packageStable},
			},
			DefaultChannelName: stableChannel,
		},
	}

	crdPlural := genName("ins")
	crdName := crdPlural + ".cluster.com"
	crd := newCRD(crdName, testNamespace, crdPlural)
	namedStrategy := newNginxInstallStrategy(genName("dep-"))
	csv := newCSV(packageStable, testNamespace, "", *semver.New("0.1.0"), []extv1beta1.CustomResourceDefinition{crd}, nil, namedStrategy)

	c := newKubeClient(t)
	crc := newCRClient(t)

	catalogSourceName := genName("mock-ocs")
	_, cleanupCatalogSource, err := createInternalCatalogSource(t, c, crc, catalogSourceName, testNamespace, manifests, []extv1beta1.CustomResourceDefinition{crd}, []v1alpha1.ClusterServiceVersion{csv})
	require.NoError(t, err)
	defer cleanupCatalogSource()

	expectedStatus := packagev1alpha1.PackageManifestStatus{
		CatalogSourceName:      catalogSourceName,
		CatalogSourceNamespace: testNamespace,
		PackageName:            packageName,
		Channels: []packagev1alpha1.PackageChannel{
			packagev1alpha1.PackageChannel{
				Name:           stableChannel,
				CurrentCSVName: packageStable,
				CurrentCSVDesc: packagev1alpha1.CreateCSVDescription(&csv),
			},
		},
		DefaultChannelName: stableChannel,
	}

	// get PackageManifest
	pmc := newPMClient(t)
	pm, err := fetchPackageManifest(t, pmc, testNamespace, packageName, packageManifestHasStatus)

	// check against expected
	require.NoError(t, err, "error getting package manifest")
	require.NotNil(t, pm)
	require.Equal(t, packageName, pm.GetName())
	require.Equal(t, expectedStatus, pm.Status)
}
