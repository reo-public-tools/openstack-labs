package theforeman

func CheckOutVLANDomain(url string, session string, domainid int) (error) {
    err := SetDomainParameter(url, session, domainid, "in-use", "yes")
    if err != nil {
        return err
    }
    return nil
}

func ReleaseVLANDomain(url string, session string, domainid int) (error) {
    err := SetDomainParameter(url, session, domainid, "in-use", "no")
    if err != nil {
        return err
    }
    return nil
}

