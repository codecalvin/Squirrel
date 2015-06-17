//
//  AppDelegate.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "GeneralTabBarViewController.h"

#define DICTIONARY_TAB_BAR_TITLE1 @"Register"
#define DICTIONARY_TAB_BAR_TITLE2 @"Publish"
#define DICTIONARY_TAB_BAR_TITLE3 @"BBS"

#define DICTIONARY_TAB_BAR_TITLE4 @"Me"
#define DICTIONARY_TAB_BAR_TITLE5 @"Me"

@interface DictionaryTabBarViewController : GeneralTabBarViewController
{

}

+ (DictionaryTabBarViewController*)singleton;

@end
