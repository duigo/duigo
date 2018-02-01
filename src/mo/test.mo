# This is comment

type  Auth
{
    exts		map(string,string)
	username	string
	password	string
	exts		map(string,type  {
                               a yyt
                           })
	tags		list(string)
}

type 	UserInfo
{
	id 			string
	height		float
	email		string
	address()	string
}

model	userinfo	UserInfo


model	userinfo2	type UserInfo2 {
    a int
    b string
    c map(string, int)
}


model	userinfo3	type  {
    a int
    b string
    c map(string, int)
}