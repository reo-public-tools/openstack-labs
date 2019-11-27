package theforeman

func CheckOutVLANDomain(url string, session string, domainid int) {
    SetDomainParameter(url, session, domainid, "in-use", "yes")
}

func ReleaseVLANDomain(url string, session string, domainid int) {
    SetDomainParameter(url, session, domainid, "in-use", "no")
}

