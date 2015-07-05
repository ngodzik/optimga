// Work in progress

package optimga

type Selecter interface {
	Select(parents *Pop, offspring *Pop)
}
