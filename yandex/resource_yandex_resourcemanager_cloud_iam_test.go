package yandex

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/access"
)

func cloudIamBindingImportStep(resourceName, cloudID, role string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateId:     fmt.Sprintf("%s %s", cloudID, role),
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM binding can be applied to a cloud
func TestAccCloudIamBinding_basic(t *testing.T) {
	cloudID := getExampleCloudID()
	role := "viewer"
	userID := getExampleUserID2()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Apply an IAM binding
			{
				Config: testAccCloudAssociateBindingBasic(cloudID, role, userID),
				Check: testAccCheckCloudIam(
					"yandex_resourcemanager_cloud_iam_binding.acceptance", role, []string{"userAccount:" + userID}),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role),
		},
	})
}

// Test that multiple IAM bindings can be applied to a cloud, one at a time
func TestAccCloudIamBinding_multiple(t *testing.T) {
	cloudID := getExampleCloudID()
	role1 := "editor"
	role2 := "viewer"
	userID1 := getExampleUserID1()
	userID2 := getExampleUserID2()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Apply an IAM binding
			{
				Config: testAccCloudAssociateBindingBasic(cloudID, role1, userID1),
				Check:  testAccCheckCloudIam("yandex_resourcemanager_cloud_iam_binding.acceptance", role1, []string{"userAccount:" + userID1}),
			},
			// Apply another IAM binding
			{
				Config: testAccCloudAssociateBindingMultiple(cloudID, role1, role2, userID1, userID2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudIam("yandex_resourcemanager_cloud_iam_binding.acceptance", role1,
						[]string{"userAccount:" + userID1, "userAccount:" + userID2}),
					testAccCheckCloudIam("yandex_resourcemanager_cloud_iam_binding.multiple", role2,
						[]string{"userAccount:" + userID1, "userAccount:" + userID2}),
				),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role1),
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.multiple", cloudID, role2),
		},
	})
}

// Test that multiple IAM bindings can be applied to a cloud all at once
func TestAccCloudIamBinding_multipleAtOnce(t *testing.T) {
	cloudID := getExampleCloudID()
	role := "editor"
	role2 := "viewer"
	userID1 := getExampleUserID1()
	userID2 := getExampleUserID2()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Prepare data source about cloud ID
			{
				Config: testAccCheckResourceManagerCloud_byID(cloudID),
			},
			// Apply an IAM binding
			{
				Config: testAccCloudAssociateBindingMultiple(cloudID, role, role2, userID1, userID2),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role),
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.multiple", cloudID, role2),
		},
	})
}

// Test that an IAM binding can be updated once applied to a cloud
func TestAccCloudIamBinding_update(t *testing.T) {
	cloudID := getExampleCloudID()
	role := "editor"
	userID1 := getExampleUserID1()
	userID2 := getExampleUserID2()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Apply an IAM binding
			{
				Config: testAccCloudAssociateBindingBasic(cloudID, role, userID1),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role),
			// Apply an updated IAM binding
			{
				Config: testAccCloudAssociateBindingUpdated(cloudID, role, userID1, userID2),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role),
			// Drop the original member
			{
				Config: testAccCloudAssociateBindingDropMemberFromBasic(cloudID, role, userID1),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role),
		},
	})
}

// Test that an IAM binding can be removed from a cloud
func TestAccCloudIamBinding_remove(t *testing.T) {
	cloudID := getExampleCloudID()
	role1 := "admin"
	role2 := "viewer"
	userID1 := getExampleUserID1()
	userID2 := getExampleUserID2()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Prepare data source about cloud ID
			{
				Config: testAccCheckResourceManagerCloud_byID(cloudID),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccCloudAssociateBindingMultiple(cloudID, role1, role2, userID1, userID2),
			},
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.acceptance", cloudID, role1),
			cloudIamBindingImportStep("yandex_resourcemanager_cloud_iam_binding.multiple", cloudID, role2),
			// Remove the bindings
			{
				Config: testAccCheckResourceManagerCloud_byID(cloudID),
			},
		},
	})
}

// Test that an IAM member can be applied to a cloud
func TestAccCloudIamMember_basic(t *testing.T) {
	cloudID := getExampleCloudID()
	role := "viewer"
	userID := getExampleUserID1()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Apply an IAM member
			{
				Config: testAccCloudAssociateMemberBasic(cloudID, role, userID),
			},
			{
				ResourceName:      "yandex_resourcemanager_cloud_iam_member.acceptance",
				ImportStateId:     fmt.Sprintf("%s %s %s", cloudID, role, "userAccount:"+userID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudIam(resourceName, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		cloudID := strings.SplitN(rs.Primary.ID, "/", 2)[0]

		resp, err := config.sdk.ResourceManager().Cloud().ListAccessBindings(context.Background(), &access.ListAccessBindingsRequest{
			ResourceId: cloudID,
		})

		if err != nil {
			return err
		}

		var roleMembers []string
		for _, binding := range resp.AccessBindings {
			if binding.RoleId == role {
				member := binding.Subject.Type + ":" + binding.Subject.Id
				roleMembers = append(roleMembers, member)
			}
		}
		sort.Strings(members)
		sort.Strings(roleMembers)

		if reflect.DeepEqual(members, roleMembers) {
			return nil
		}

		return fmt.Errorf("Binding found but expected members is %v, got %v", members, roleMembers)
	}
}

func testAccCloudAssociateBindingBasic(cloudID, role, userID string) string {
	prerequisiteMembership, deps := testAccCloudAssignCloudMemberRole(cloudID, userID)
	return prerequisiteMembership + fmt.Sprintf(`
resource "yandex_resourcemanager_cloud_iam_binding" "acceptance" {
  cloud_id = "%s"
  role     = "%s"
  members  = ["userAccount:%s"]

  depends_on = [%s]
}
`, cloudID, role, userID, deps)
}

func testAccCloudAssociateBindingMultiple(cloudID, role1, role2, userID1, userID2 string) string {
	prerequisiteMembership, deps := testAccCloudAssignCloudMemberRole(cloudID, userID1, userID2)

	multiple1 := fmt.Sprintf(`
resource "yandex_resourcemanager_cloud_iam_binding" "acceptance" {
  cloud_id = "%s"
  role     = "%s"
  members  = ["userAccount:%s", "userAccount:%s"]

  depends_on = [%s]
}
`, cloudID, role1, userID1, userID2, deps)

	multiple2 := fmt.Sprintf(`
resource "yandex_resourcemanager_cloud_iam_binding" "multiple" {
  cloud_id = "%s"
  role     = "%s"
  members  = ["userAccount:%s", "userAccount:%s"]

  depends_on = [%s]
}
`, cloudID, role2, userID1, userID2, deps)

	return prerequisiteMembership + multiple1 + multiple2
}

func testAccCloudAssociateBindingUpdated(cloudID, role, userID1, userID2 string) string {
	prerequisiteMembership, deps := testAccCloudAssignCloudMemberRole(cloudID, userID1, userID2)
	return prerequisiteMembership + fmt.Sprintf(`
resource "yandex_resourcemanager_cloud_iam_binding" "acceptance" {
  cloud_id = "%s"
  role     = "%s"
  members  = ["userAccount:%s", "userAccount:%s"]

  depends_on = [%s]
}
`, cloudID, role, userID1, userID2, deps)
}

func testAccCloudAssociateBindingDropMemberFromBasic(cloudID, role, userID string) string {
	prerequisiteMembership, deps := testAccCloudAssignCloudMemberRole(cloudID, userID)
	return prerequisiteMembership + fmt.Sprintf(`
resource "yandex_resourcemanager_cloud_iam_binding" "acceptance" {
  cloud_id = "%s"
  role     = "%s"
  members  = ["userAccount:%s"]

  depends_on = [%s]
}
`, cloudID, role, userID, deps)
}

func testAccCloudAssociateMemberBasic(cloudID, role, userID string) string {
	prerequisiteMembership, deps := testAccCloudAssignCloudMemberRole(cloudID, userID)

	return prerequisiteMembership + fmt.Sprintf(`
resource "yandex_resourcemanager_cloud_iam_member" "acceptance" {
  cloud_id = "%s"
  role     = "%s"
  member   = "userAccount:%s"

  depends_on = [%s]
}
`, cloudID, role, userID, deps)
}

func testAccCloudAssignCloudMemberRole(cloudID string, usersID ...string) (string, string) {
	var config bytes.Buffer
	var resourceRefs []string

	for _, userID := range usersID {
		resType := "yandex_resourcemanager_cloud_iam_member"
		resName := fmt.Sprintf("membership-%s-%s", cloudID, userID)

		config.WriteString(fmt.Sprintf(`
// Make user member of cloud to allow assign another roles
resource "%s" "%s" {
  cloud_id = "%s"
  role     = "resource-manager.clouds.member"
  member   = "userAccount:%s"
}
`, resType, resName, cloudID, userID))

		resourceRefs = append(resourceRefs, fmt.Sprintf("\"%s.%s\"", resType, resName))
	}

	return config.String(), strings.Join(resourceRefs, ",")
}
