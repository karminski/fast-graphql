-- jit.lua
_G.result  = ""
_G.User    = {}
_G.Friends = {}

-- resolve SelectionSet
function m0()
    -- init
    local buf = {}
    local queryVariablesMap = {} -- generate from go runtime
    
    -- build data header
    buf[#buf+1] = "{\"data\":"
    
    -- get user data by param from [go runtime]
    local user = _G.User 
    
    -- resolve user
    resolveUser(user)

    -- build no error data end
    buf[#buf+1] = ",\"errors\":null,\"jit-result\":true}"

    -- finally, return result to go runtime
    local r = table.concat(buf)
    _G.result = r
end

function resolveUser(user)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"Id\":"
    buf[#buf+1] = user.Id
    buf[#buf+1] = ",\"Name\":"
    buf[#buf+1] = user.Name
    buf[#buf+1] = ",\"Email\":"
    buf[#buf+1] = user.Email
    buf[#buf+1] = ",\"Married\":"
    buf[#buf+1] = user.Married
    buf[#buf+1] = ",\"Height\":"
    buf[#buf+1] = user.Height
    buf[#buf+1] = ",\"Gender\":"
    buf[#buf+1] = user.Gender
    buf[#buf+1] = ",\"Friends\":"
    resolveFriends(user.Friends)
    buf[#buf+1] = ",\"Location\":"
    resolveLocation(user.Location)
    buf[#buf+1] = "}"
end

function resolveFriends(friendIDs)
    -- get friends details from [go runtime]

    local friends = _G.Friends

    -- resolve friends
    l = #friends
    buf[#buf+1] = "["
    for i=1, l do
        resolveFriend(friends[i])
        if i <> l then
            buf[#buf+1] = ","
        end
    end
    buf[#buf+1] = "]"
end

function resolveFriend(friend)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"Id\":"
    buf[#buf+1] = friend.Id
    buf[#buf+1] = ",\"Name\":"
    buf[#buf+1] = friend.Name
    buf[#buf+1] = ",\"Email\":"
    buf[#buf+1] = friend.Email
    buf[#buf+1] = ",\"Married\":"
    buf[#buf+1] = friend.Married
    buf[#buf+1] = ",\"Height\":"
    buf[#buf+1] = friend.Height
    buf[#buf+1] = ",\"Gender\":"
    buf[#buf+1] = friend.Gender
    buf[#buf+1] = ",\"Location\":"
    resolveLocation(friend.Location)
    buf[#buf+1] = "}"
end

function resolveLocation(location)
    buf[#buf+1] = "{"
    buf[#buf+1] = "\"City\":"
    buf[#buf+1] = location.City
    buf[#buf+1] = ",\"Country\":"
    buf[#buf+1] = location.Country
    buf[#buf+1] = "}"
end


