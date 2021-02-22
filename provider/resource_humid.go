package provider

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kscarlett/humid"
	"github.com/kscarlett/humid/wordlist"
)

func generate() *schema.Resource {
	return &schema.Resource{
		Description: `
The resource ` + "`random_id`" + ` generates random IDs that are intended to be
used as unique identifiers for other resources.
This resource *does* use a cryptographic random number generator in order
to minimize the chance of collisions, but the space for random IDs strongly
depends on the wordlists and other settings used.
This resource can be used in conjunction with resources that have
the ` + "`create_before_destroy`" + ` lifecycle flag set to avoid conflicts with
unique names during the brief period where both the old and new resources
exist concurrently.
		`,
		Create: CreateHumid,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportHumid,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Description: "Arbitrary map of values that, when changed, will trigger recreation of " +
					"resource. See [the official random provider documentation](https://registry.terraform.io/providers/hashicorp/random/latest/docs#resource-keepers) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"wordlist": {
				Description: "Name of the wordlist to use to generate the random ID." +
					"See the [humid repo](https://github.com/kscarlett/humid/tree/main/wordlist) for more information on the wordlists available." +
					"Defaults to the \"animals\" list.",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"adjectives": {
				Description: "Amount of adjectives to use to generate the ID. Adds a lot more options." +
					"Defaults to 1.",
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"separator": {
				Description: "What to use between words in the ID. Defaults to \"-\".",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"capitalize": {
				Description: "Whether to capitalize the first letter of each word or to leave everything lowercase." +
					"Defaults to false (lowercase).",
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"id": {
				Description: "The generated ID presented in string format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"result": {
				Description: "The generated ID presented in string format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func CreateHumid(d *schema.ResourceData, meta interface{}) error {
	// get list from d.Get(wordlist)
	listName, listSet := d.GetOk("wordlist")
	if !listSet {
		listName = "animals"
	}

	list, err := getList(listName.(string))
	if err != nil {
		return err
	}

	adjCount, adjCountSet := d.GetOk("adjectives")
	if !adjCountSet {
		adjCount = 1
	}

	sep, sepSet := d.GetOk("separator")
	if !sepSet {
		sep = "-"
	}

	cap, capSet := d.GetOk("capitalize")
	if !capSet {
		cap = false
	}

	id := humid.GenerateWithOptions(&humid.Options{
		List:           list,
		AdjectiveCount: adjCount.(int),
		Separator:      sep.(string),
		Capitalize:     cap.(bool),
	})

	d.Set("result", id)
	d.SetId(id)

	return nil
}

func ImportHumid(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	d.Set("result", id)
	d.SetId(id)

	// not sure how to import everything else
	return []*schema.ResourceData{d}, nil
}

// Warning, this now needs to be manually updated whenever the Humid package releases a new list!!
func getList(s string) (list []string, err error) {
	switch strings.ToLower(s) {
	case "animals":
		return wordlist.Animals, nil
	case "adjectives":
		return wordlist.Adjectives, nil
	default:
		return nil, errors.New("wordlist " + s + " could not be found in the humid package")
	}
}
